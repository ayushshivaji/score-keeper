export interface User {
  id: string;
  google_id: string;
  email: string;
  name: string;
  avatar_url: string | null;
  matches_played: number;
  matches_won: number;
  total_points: number;
  created_at: string;
  updated_at: string;
}

export interface UserProfile extends User {
  losses: number;
  win_rate: number;
  current_streak: number;
  longest_win_streak: number;
  longest_loss_streak: number;
  recent_form: string[];
}
