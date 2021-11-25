

func (r *Reversi) matching(ctx context.Context, cli pb.MatchingServiceClient) error {
	// マッチングリクエスト
	stream, err := cli.JoinRoom(ctx, &pb.JoinRoomRequest{})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	fmt.Println("Requested matching...")

	// ストリーミングでレスポンスを受け取る
	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		if resp.GetStatus() == pb.JoinRoomResponse_MATCHED {
			// マッチング成立
			r.room = build.Room(resp.GetRoom())
			r.me = bulid.Player(resp.GetMe())
			fmt.Printf("Matched room_id = %d\n", resp.GetRoom().GetId())
			return nil
		} else if resp.GetStatus() == pb.JoinRoomResponse_WAITING {
			// 待機中
			fmt.Println("Waiting matching...")
		}
	}
}

func (r *Reversi) play(ctx context.Context, cli pb.GameServiceClient) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	// 双方向ストリーミングを開始する
	stream, err := cli.Play(c)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	go func() {
		err := r.send(c, stream)
		if err != nil {
			cancel()
		}
	}()

	err = r.receive(c, stream)
	if err != nil {
		cancel()
		return err
	}

	return nil
}

func (r *Reversi) receive(ctx context.Context, stream pb.GameService_PlayClient) error {
	for {
		// サーバーからのストリーミングを受け取る
		res, err := stream.Recv()

		switch res.GetEvent().(type) {
		case *pb.PlayResponse_Waiting:
			// 開始待機中
		case *pb.PlayResponse_Ready:
			// 開始
			r.started = true
			r.game.Display(r.me.Color)
		case *pb.PlayResponse_Move:
			// 手を打たれた
			color := build.Color(res.GetMove().GetPlayer().GetColor())
			if color != r.me.Color {
				move := res.GetMove().GetMove()
				// クライアント側のゲーム情報に反映させる
				r.game.Move(move.GetX(), move.GetY(), color)
			}
		case *pb.PlayResponse_Finished:
			// ゲームが終了した
			r.finished = true

			return nil
		}

		select {
		case <-ctx.Done():
			// キャンセルされたので終了する
			return nil
		default:
		}
	}
}

func (r *Reversi) send(ctx context.Context, stream pb.GameService_PlayClient) error {
	for {
		if r.finished {
			return nil
		} else if !r.started {
			// 未開始なので開始リクエストを送る
			err := stream.Send(&pb.PlayRequest{
				RoomId: r.room.ID,
				Player: build.PBPlayer(r.me),
				Action: &pb.PlayRequest_Start{
					Start: &pb.PlayRequest_StartAction{},
				},
			})

			for {
				// 相手が開始するまで待機する
				if r.started {
					// 開始をreceiveした
					fmt.Println("READY GO!")
					break
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			// 手の入力を待機する
			fmt.Print("Input Your Move (ex. A-1): ")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()

			// 入力された手を解析する
			text := stdin.Text()
			x, y, err := parseInput(text)

			// 手を打つ
			_, err = r.game.Move(x, y, r.me.Color)

			go func() {
				// サーバーに手を送る
				err = stream.Send(&pb.PlayRequest{
					RoomId: r.room.ID,
					Player: build.PBPlayer(r.me),
					Action: &pb.PlayRequest_Move{
						Move: &pb.PlayRequest_MoveAction{
							Move: &pb.Move{
								X: x,
								Y: y,
							},
						},
					},
				})
			}()

			// 一度手を打ったら5秒間待機する
			ch := make(chan int)
			go func(ch chan int) {
				for i := 0; i < 5; i++ {
					time.Sleep(1 * time.Second)
				}
				ch <- 0
			}(ch)
			<-ch
		}

		select {
		case <-ctx.Done():
			// キャンセルされたので終了する
			return nil
		default:
		}
	}
}
