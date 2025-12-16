'use client';

import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useUpdateMe } from '@/api/public/user/user';
import { useAuth } from '@/contexts/auth-context';

interface ProfileEditFormData {
  name: string;
}

interface ProfileEditDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function ProfileEditDialog({ open, onOpenChange }: ProfileEditDialogProps) {
  const { user, refreshUser } = useAuth();
  const updateMe = useUpdateMe();

  const { register, handleSubmit, reset, formState: { errors } } = useForm<ProfileEditFormData>({
    defaultValues: {
      name: user?.name ?? '',
    },
  });

  useEffect(() => {
    if (open && user) {
      reset({
        name: user.name ?? '',
      });
    }
  }, [open, user, reset]);

  const onSubmit = async (data: ProfileEditFormData) => {
    try {
      await updateMe.mutateAsync({
        data: {
          name: data.name,
        },
      });
      await refreshUser();
      toast.success('プロフィールを更新しました');
      onOpenChange(false);
    } catch (error) {
      toast.error('プロフィールの更新に失敗しました');
      console.error(error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>プロフィール編集</DialogTitle>
            <DialogDescription>
              表示名を変更できます。
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="profile-name">名前</Label>
              <Input
                id="profile-name"
                placeholder="名前を入力"
                {...register('name', { required: '名前は必須です' })}
              />
              {errors.name && (
                <p className="text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label className="text-muted-foreground">メールアドレス</Label>
              <p className="text-sm">{user?.email}</p>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              キャンセル
            </Button>
            <Button type="submit" disabled={updateMe.isPending}>
              {updateMe.isPending ? '更新中...' : '更新'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
