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
  /** 作成者カラムを表示するか（チーム公開タブ用） */
  showCreator?: boolean;
}

export function TodoList({
  todos,
  emptyMessage = 'Todoがありません',
  allEditable = false,
  myTodoIds,
  showCreator = false,
}: TodoListProps) {
  if (todos.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500 bg-white rounded-lg border border-gray-200">
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden overflow-x-auto">
      <table className="w-full table-fixed min-w-[600px]">
        <colgroup>
          <col className="w-14" />
          <col />
          {showCreator && <col className="w-28" />}
          <col className="w-28" />
          <col className="w-24" />
          <col className="w-28" />
        </colgroup>
        <thead>
          <tr className="bg-gray-50 border-b border-gray-200">
            <th className="py-3 px-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              完了
            </th>
            <th className="py-3 px-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              タイトル
            </th>
            {showCreator && (
              <th className="py-3 px-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                作成者
              </th>
            )}
            <th className="py-3 px-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              期限
            </th>
            <th className="py-3 px-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
              公開
            </th>
            <th className="py-3 px-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              操作
            </th>
          </tr>
        </thead>
        <tbody>
          {todos.map((todo) => {
            const isOwner = allEditable || (myTodoIds !== undefined && todo.id !== undefined && myTodoIds.has(todo.id));
            return (
              <TodoItem
                key={todo.id}
                todo={todo}
                isOwner={isOwner}
                showCreator={showCreator}
              />
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
