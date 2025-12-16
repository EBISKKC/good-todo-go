'use client';

import { useState } from 'react';
import { Checkbox } from '@/components/ui/checkbox';
import { Button } from '@/components/ui/button';
import { Trash2, Globe, Lock, Pencil, Eye } from 'lucide-react';
import { TodoResponse } from '@/api/public/model/components-schemas-todo';
import { useUpdateTodo, useDeleteTodo, getGetTodosQueryKey, getGetPublicTodosQueryKey } from '@/api/public/todo/todo';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { format } from 'date-fns';
import { ja } from 'date-fns/locale';
import { EditTodoDialog } from './edit-todo-dialog';
import { ViewTodoDialog } from './view-todo-dialog';

interface TodoItemProps {
  todo: TodoResponse;
  isOwner?: boolean;
  showCreator?: boolean;
}

export function TodoItem({ todo, isOwner = false, showCreator = false }: TodoItemProps) {
  const queryClient = useQueryClient();
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isViewDialogOpen, setIsViewDialogOpen] = useState(false);

  const updateTodo = useUpdateTodo();
  const deleteTodo = useDeleteTodo();

  const handleToggleComplete = async () => {
    if (!todo.id || !isOwner) return;

    try {
      await updateTodo.mutateAsync({
        todoId: todo.id,
        data: { completed: !todo.completed },
      });
      queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() });
      queryClient.invalidateQueries({ queryKey: getGetPublicTodosQueryKey() });
    } catch (error) {
      toast.error('更新に失敗しました');
      console.error(error);
    }
  };

  const handleTogglePublic = async () => {
    if (!todo.id || !isOwner) return;

    try {
      await updateTodo.mutateAsync({
        todoId: todo.id,
        data: { is_public: !todo.is_public },
      });
      queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() });
      queryClient.invalidateQueries({ queryKey: getGetPublicTodosQueryKey() });
      toast.success(todo.is_public ? '非公開にしました' : '公開しました');
    } catch (error) {
      toast.error('更新に失敗しました');
      console.error(error);
    }
  };

  const handleDelete = async () => {
    if (!todo.id || !isOwner) return;

    setIsDeleting(true);
    try {
      await deleteTodo.mutateAsync({ todoId: todo.id });
      queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() });
      queryClient.invalidateQueries({ queryKey: getGetPublicTodosQueryKey() });
      toast.success('削除しました');
    } catch (error) {
      toast.error('削除に失敗しました');
      console.error(error);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <>
      <tr className={`border-b border-gray-100 hover:bg-gray-50 transition-colors ${todo.completed ? 'opacity-50' : ''}`}>
        {/* チェックボックス */}
        <td className="py-3 px-3">
          {isOwner ? (
            <Checkbox
              checked={todo.completed ?? false}
              onCheckedChange={handleToggleComplete}
              disabled={updateTodo.isPending}
            />
          ) : (
            <div className="w-4 h-4 rounded border border-gray-300 flex items-center justify-center">
              {todo.completed && <div className="w-2 h-2 bg-gray-400 rounded-sm" />}
            </div>
          )}
        </td>

        {/* タイトル・説明 */}
        <td className="py-3 px-3">
          <div className="min-w-0 overflow-hidden">
            <p className={`font-medium text-gray-900 truncate ${todo.completed ? 'line-through text-gray-500' : ''}`}>
              {todo.title}
            </p>
            {todo.description && (
              <p className="text-sm text-gray-500 truncate mt-0.5">{todo.description}</p>
            )}
          </div>
        </td>

        {/* 作成者（チーム公開タブのみ） */}
        {showCreator && (
          <td className="py-3 px-3">
            {todo.created_by ? (
              <span className="text-sm text-gray-600 truncate block">
                {todo.created_by.name}
              </span>
            ) : (
              <span className="text-sm text-gray-400">-</span>
            )}
          </td>
        )}

        {/* 期限 */}
        <td className="py-3 px-3">
          {todo.due_date ? (
            <span className="text-sm text-gray-600 whitespace-nowrap">
              {format(new Date(todo.due_date), 'M/d (E)', { locale: ja })}
            </span>
          ) : (
            <span className="text-sm text-gray-400">-</span>
          )}
        </td>

        {/* 公開状態 */}
        <td className="py-3 px-3 text-center">
          <span
            className={`inline-flex items-center justify-center gap-1 text-xs px-2 py-1 rounded-full whitespace-nowrap ${
              todo.is_public
                ? 'bg-blue-50 text-blue-600'
                : 'bg-gray-100 text-gray-500'
            }`}
          >
            {todo.is_public ? (
              <Globe className="h-3 w-3" />
            ) : (
              <Lock className="h-3 w-3" />
            )}
            {todo.is_public ? '公開' : '非公開'}
          </span>
        </td>

        {/* アクション */}
        <td className="py-3 px-3">
          <div className="flex items-center justify-end gap-0.5">
            {isOwner ? (
              <>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7"
                  onClick={() => setIsEditDialogOpen(true)}
                  title="編集"
                >
                  <Pencil className="h-3.5 w-3.5 text-gray-500" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7"
                  onClick={handleTogglePublic}
                  disabled={updateTodo.isPending}
                  title={todo.is_public ? '非公開にする' : '公開する'}
                >
                  {todo.is_public ? (
                    <Lock className="h-3.5 w-3.5 text-gray-500" />
                  ) : (
                    <Globe className="h-3.5 w-3.5 text-gray-500" />
                  )}
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 text-red-500 hover:text-red-700 hover:bg-red-50"
                  onClick={handleDelete}
                  disabled={isDeleting}
                  title="削除"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </>
            ) : (
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7"
                onClick={() => setIsViewDialogOpen(true)}
                title="詳細を見る"
              >
                <Eye className="h-3.5 w-3.5 text-gray-500" />
              </Button>
            )}
          </div>
        </td>
      </tr>

      <EditTodoDialog
        todo={todo}
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
      />

      <ViewTodoDialog
        todo={todo}
        open={isViewDialogOpen}
        onOpenChange={setIsViewDialogOpen}
      />
    </>
  );
}
