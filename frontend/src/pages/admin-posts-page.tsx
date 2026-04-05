import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import type { PostListResponse } from "@/entities/post/model";
import { useSession } from "@/app/providers";
import { AdminNav } from "@/features/admin/admin-nav";
import { api } from "@/shared/api/client";
import { formatDate } from "@/shared/lib/format";
import { Card, PageHeader, SecondaryButton } from "@/shared/ui/primitives";

export function AdminPostsPage() {
  const queryClient = useQueryClient();
  const { selectedHouseId } = useSession();

  const postsQuery = useQuery({
    queryKey: ["admin-posts", selectedHouseId],
    queryFn: () => api.get<PostListResponse>(`/houses/${selectedHouseId}/posts?pageSize=50`),
    enabled: Boolean(selectedHouseId),
  });

  const deleteMutation = useMutation({
    mutationFn: (postId: number) => api.delete(`/houses/${selectedHouseId}/posts/${postId}`),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["admin-posts", selectedHouseId] });
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Admin"
        title="Управление постами"
        description="Администратор дома может удалить любой пост, даже если он создан другим пользователем."
        actions={<AdminNav />}
      />
      <div className="grid gap-4">
        {postsQuery.data?.items.map((post) => (
          <Card key={post.id} className="space-y-4">
            <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
              <div className="space-y-2">
                <div className="text-xs uppercase tracking-[0.2em] text-slate-500">{post.categoryName}</div>
                <h3 className="text-xl font-semibold text-ink">{post.title}</h3>
                <p className="text-sm leading-6 text-slate-600">{post.content}</p>
                <div className="text-sm text-slate-500">
                  Автор: {post.authorName} · {formatDate(post.createdAt)}
                </div>
              </div>
              <SecondaryButton onClick={() => deleteMutation.mutate(post.id)}>Удалить пост</SecondaryButton>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
