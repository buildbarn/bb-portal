export class NotFoundError extends Error {
  type?: string;
  details?: string;
  constructor(type?: string, details?: string) {
    super("NOT_FOUND");
    this.type = type;
    this.details = details;
  }
}
