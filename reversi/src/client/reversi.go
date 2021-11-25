

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