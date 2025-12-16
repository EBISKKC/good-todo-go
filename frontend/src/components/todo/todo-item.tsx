'use client';

import { useState } from 'react';
import { Checkbox } from '@/components/ui/checkbox';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Trash2, Globe, Lock } from 'lucide-react';
import { TodoResponse } from '@/api/public/model/components-schemas-todo';
import { useUpdateTodo, useDeleteTodo, getGetTodosQueryKey, getGetPublicTodosQueryKey } from '@/api/public/todo/todo';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { format } from 'date-fns';

interface TodoItemProps {
  todo: TodoResponse;
  isOwner?: boolean;
}

export function TodoItem({ todo, isOwner = true }: TodoItemProps) {
  const queryClient = useQueryClient();
  const [isDeleting, setIsDeleting] = useState(false);

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
      toast.error('Failed to update todo');
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
      toast.success(todo.is_public ? 'Todo is now private' : 'Todo is now public');
    } catch (error) {
      toast.error('Failed to update todo');
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
      toast.success('Todo deleted');
    } catch (error) {
      toast.error('Failed to delete todo');
      console.error(error);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <Card className={`${todo.completed ? 'opacity-60' : ''}`}>
      <CardContent className="flex items-center gap-4 p-4">
        {isOwner && (
          <Checkbox
            checked={todo.completed ?? false}
            onCheckedChange={handleToggleComplete}
            disabled={updateTodo.isPending}
          />
        )}

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className={`font-medium truncate ${todo.completed ? 'line-through text-gray-500' : ''}`}>
              {todo.title}
            </h3>
            {todo.is_public ? (
              <Globe className="h-4 w-4 text-blue-500" />
            ) : (
              <Lock className="h-4 w-4 text-gray-400" />
            )}
          </div>
          {todo.description && (
            <p className="text-sm text-gray-500 truncate">{todo.description}</p>
          )}
          {todo.due_date && (
            <p className="text-xs text-gray-400 mt-1">
              Due: {format(new Date(todo.due_date), 'MMM d, yyyy')}
            </p>
          )}
        </div>

        {isOwner && (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              onClick={handleTogglePublic}
              disabled={updateTodo.isPending}
              title={todo.is_public ? 'Make private' : 'Make public'}
            >
              {todo.is_public ? (
                <Lock className="h-4 w-4" />
              ) : (
                <Globe className="h-4 w-4" />
              )}
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleDelete}
              disabled={isDeleting}
              className="text-red-500 hover:text-red-700"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
