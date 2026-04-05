export interface House {
  id: number;
  name: string;
  address: string;
  createdBy: number;
  createdAt: string;
  role: "resident" | "admin";
}
