'use client';

import { TodoItem } from './todo-item';
import { TodoResponse } from '@/api/public/model/components-schemas-todo';

interface TodoListProps {
  todos: TodoResponse[];
  isOwner?: boolean;
  emptyMessage?: string;
}

export function TodoList({ todos, isOwner = true, emptyMessage = 'No todos yet' }: TodoListProps) {
  if (todos.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {todos.map((todo) => (
        <TodoItem key={todo.id} todo={todo} isOwner={isOwner} />
      ))}
    </div>
  );
}
