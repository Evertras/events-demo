// IOpenGame represents a game that has been created, but requires
// more players to start.
export interface IOpenGame {
  id: number;
  name: string;
  created: Date;
}
