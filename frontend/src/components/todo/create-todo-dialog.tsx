'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { useCreateTodo, getGetTodosQueryKey, getGetPublicTodosQueryKey } from '@/api/public/todo/todo';

interface CreateTodoFormData {
  title: string;
  description?: string;
  is_public: boolean;
  due_date?: string;
}

export function CreateTodoDialog() {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const { register, handleSubmit, reset, formState: { errors }, watch, setValue } = useForm<CreateTodoFormData>({
    defaultValues: {
      is_public: false,
    },
  });

  const createTodo = useCreateTodo();
  const isPublic = watch('is_public');

  const onSubmit = async (data: CreateTodoFormData) => {
    try {
      await createTodo.mutateAsync({
        data: {
          title: data.title,
          description: data.description,
          is_public: data.is_public,
          due_date: data.due_date ? new Date(data.due_date).toISOString() : undefined,
        },
      });
      queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() });
      queryClient.invalidateQueries({ queryKey: getGetPublicTodosQueryKey() });
      toast.success('Todo created successfully');
      reset();
      setOpen(false);
    } catch (error) {
      toast.error('Failed to create todo');
      console.error(error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          New Todo
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Create New Todo</DialogTitle>
            <DialogDescription>
              Add a new todo item to your list.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                placeholder="Enter todo title"
                {...register('title', { required: 'Title is required' })}
              />
              {errors.title && (
                <p className="text-sm text-red-500">{errors.title.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description (optional)</Label>
              <Input
                id="description"
                placeholder="Enter description"
                {...register('description')}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="due_date">Due Date (optional)</Label>
              <Input
                id="due_date"
                type="date"
                {...register('due_date')}
              />
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox
                id="is_public"
                checked={isPublic}
                onCheckedChange={(checked) => setValue('is_public', checked === true)}
              />
              <Label htmlFor="is_public" className="text-sm font-normal">
                Make this todo public (visible to your team)
              </Label>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createTodo.isPending}>
              {createTodo.isPending ? 'Creating...' : 'Create'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
