export interface Post {
  id: number;
  houseId: number;
  authorId: number;
  categoryId: number;
  title: string;
  content: string;
  imageUrl?: string;
  commentsCount: number;
  createdAt: string;
  updatedAt: string;
  authorName: string;
  categoryName: string;
}

export interface PostListResponse {
  items: Post[];
  total: number;
  page: number;
  pageSize: number;
}
