'use client';

import { TodoItem } from './todo-item';
import { TodoResponse } from '@/api/public/model/components-schemas-todo';

interface TodoListProps {
  todos: TodoResponse[];
  emptyMessage?: string;
  /** trueの場合、すべてのTodoが編集可能（マイTodoタブ用） */
  allEditable?: boolean;
  /** 自分のTodoのIDセット（チーム公開Todoタブで自分のTodoを編集可能にするため） */
  myTodoIds?: Set<string>;
}

export function TodoList({
  todos,
  emptyMessage = 'Todoがありません',
  allEditable = false,
  myTodoIds,
}: TodoListProps) {
  if (todos.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {todos.map((todo) => {
        const isOwner = allEditable || (myTodoIds !== undefined && todo.id !== undefined && myTodoIds.has(todo.id));
        return (
          <TodoItem
            key={todo.id}
            todo={todo}
            isOwner={isOwner}
          />
        );
      })}
    </div>
  );
}
