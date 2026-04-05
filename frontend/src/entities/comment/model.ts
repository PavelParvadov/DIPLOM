export interface Comment {
  id: number;
  postId: number;
  authorId: number;
  content: string;
  createdAt: string;
  authorName: string;
}

export interface CommentListResponse {
  items: Comment[];
  total: number;
  page: number;
  pageSize: number;
}
