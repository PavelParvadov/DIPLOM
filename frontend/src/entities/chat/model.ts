export interface ChatMessage {
  id: number;
  houseId: number;
  authorId: number;
  content: string;
  imageUrl?: string;
  createdAt: string;
  authorName: string;
}
