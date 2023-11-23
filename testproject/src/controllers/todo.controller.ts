export class TodoController {
  constructor(private readonly todoService: TodoService) {}

  @Get()
  async handler() {
    return null;
  }
}
