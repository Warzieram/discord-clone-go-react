export type Message = {
  id?: number;
  sender: string;
  created_at: string;
  content: string;
  room_id?: number;
};

export type BroadcastedMessage = {
  command_type: string;
  data: Message | number;
};
