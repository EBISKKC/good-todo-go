'use client';

import { useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { LogOut, User, Building2 } from 'lucide-react';
import { useAuth } from '@/contexts/auth-context';
import { useGetTodos, useGetPublicTodos } from '@/api/public/todo/todo';
import { TodoList } from '@/components/todo/todo-list';
import { CreateTodoDialog } from '@/components/todo/create-todo-dialog';

export default function TodosPage() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading: authLoading, logout } = useAuth();

  const { data: myTodosData, isLoading: myTodosLoading } = useGetTodos(
    undefined,
    { query: { enabled: isAuthenticated } }
  );

  const { data: publicTodosData, isLoading: publicTodosLoading } = useGetPublicTodos(
    undefined,
    { query: { enabled: isAuthenticated } }
  );

  // マイTodoのIDセットを作成（チーム公開Todoで自分のTodoを編集可能にするため）
  const myTodoIds = useMemo(() => {
    const todos = myTodosData?.todos ?? [];
    return new Set(todos.map((todo) => todo.id).filter((id): id is string => id !== undefined));
  }, [myTodosData?.todos]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [authLoading, isAuthenticated, router]);

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  const myTodos = myTodosData?.todos ?? [];
  const publicTodos = publicTodosData?.todos ?? [];

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-bold">Good Todo</h1>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-3 text-sm">
              <div className="flex items-center gap-1.5 text-gray-500 bg-gray-100 px-2 py-1 rounded">
                <Building2 className="h-3.5 w-3.5" />
                <span>{user?.tenantSlug}</span>
              </div>
              <div className="flex items-center gap-1.5 text-gray-600">
                <User className="h-4 w-4" />
                <span>{user?.name}</span>
              </div>
            </div>
            <Button variant="ghost" size="sm" onClick={logout}>
              <LogOut className="h-4 w-4 mr-2" />
              ログアウト
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold">Todo一覧</h2>
          <CreateTodoDialog />
        </div>

        <Tabs defaultValue="my-todos" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="my-todos">
              マイTodo ({myTodos.length})
            </TabsTrigger>
            <TabsTrigger value="public-todos">
              チーム公開Todo ({publicTodos.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="my-todos" className="mt-6">
            {myTodosLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
              </div>
            ) : (
              <TodoList
                todos={myTodos}
                allEditable={true}
                emptyMessage="Todoがありません。最初のTodoを作成しましょう！"
              />
            )}
          </TabsContent>

          <TabsContent value="public-todos" className="mt-6">
            {publicTodosLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
              </div>
            ) : (
              <TodoList
                todos={publicTodos}
                myTodoIds={myTodoIds}
                emptyMessage="チームで公開されているTodoはまだありません。"
              />
            )}
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
