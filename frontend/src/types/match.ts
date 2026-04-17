import { User } from "./user";

export interface Match {
  id: string;
  player1_id: string;
  player2_id: string;
  winner_id: string;
  player1_score: number;
  player2_score: number;
  played_at: string;
  created_at: string;
  created_by: string;
  player1: User;
  player2: User;
}

export interface CreateMatchRequest {
  player1_id: string;
  player2_id: string;
  player1_score: number;
  player2_score: number;
  played_at: string;
}

export interface HeadToHead {
  player1: User;
  player2: User;
  total_matches: number;
  player1_wins: number;
  player2_wins: number;
  player1_points: number;
  player2_points: number;
  matches: Match[];
}
