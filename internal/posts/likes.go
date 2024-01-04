package posts

func (s Service) LikePost(userId, postId int64) error {
	post, err := s.storer.GetPost(postId)
	if err != nil {
		return err
	}

	if err := s.storer.InsertPostLike(userId, post.ID); err != nil {
		return err
	}
	return nil
}

func (s Service) DislikePost(userId, postId int64) error {
	post, err := s.storer.GetPost(postId)
	if err != nil {
		return err
	}

	if err := s.storer.DeletePostLike(userId, post.ID); err != nil {
		return err
	}
	return nil
}
