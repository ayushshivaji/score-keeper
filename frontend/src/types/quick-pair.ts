import { User } from "@/types/user";

export interface QuickPair {
  id: string;
  user_id: string;
  player1_id: string;
  player2_id: string;
  created_at: string;
  player1: User;
  player2: User;
}
