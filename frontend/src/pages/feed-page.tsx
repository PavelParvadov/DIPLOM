import { useQuery } from "@tanstack/react-query";
import { Link, useSearchParams } from "react-router-dom";

import { useSession } from "@/app/providers";
import type { Category } from "@/entities/category/model";
import type { PostListResponse } from "@/entities/post/model";
import { PostCard } from "@/features/posts/post-card";
import { api } from "@/shared/api/client";
import { EmptyState, PageHeader, Select, SecondaryButton } from "@/shared/ui/primitives";

export function FeedPage() {
  const { selectedHouseId } = useSession();
  const [searchParams, setSearchParams] = useSearchParams();
  const categoryQuery = useQuery({
    queryKey: ["categories", selectedHouseId],
    queryFn: async () => {
      const response = await api.get<{ items: Category[] }>(`/houses/${selectedHouseId}/categories`);
      return response.items;
    },
    enabled: Boolean(selectedHouseId),
  });

  const selectedCategory = searchParams.get("categoryId") ?? "";
  const postsQuery = useQuery({
    queryKey: ["posts", selectedHouseId, selectedCategory],
    queryFn: () =>
      api.get<PostListResponse>(
        `/houses/${selectedHouseId}/posts${selectedCategory ? `?categoryId=${selectedCategory}` : ""}`,
      ),
    enabled: Boolean(selectedHouseId),
  });

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Лента дома"
        title="Новости соседей"
        description="Публикации идут вертикально и листаются одна за другой, чтобы внимание было на одном посте, а не на сетке карточек."
        actions={
          <Link
            className="inline-flex items-center justify-center rounded-2xl bg-ember px-4 py-2.5 text-sm font-semibold text-white transition hover:translate-y-[-1px] hover:bg-[#c35d31]"
            to="/posts/new"
          >
            Новый пост
          </Link>
        }
      />

      <div className="flex items-center justify-end gap-3">
        <span className="rounded-full bg-white px-4 py-2 text-sm font-semibold text-ink shadow-soft">Фильтр</span>
        <Select
          className="w-[220px] bg-white shadow-soft"
          value={selectedCategory}
          onChange={(event) => {
            const next = event.target.value;
            if (next) {
              setSearchParams({ categoryId: next });
            } else {
              setSearchParams({});
            }
          }}
        >
          <option value="">Все категории</option>
          {categoryQuery.data?.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </Select>
        {selectedCategory ? (
          <SecondaryButton onClick={() => setSearchParams({})}>Сбросить</SecondaryButton>
        ) : null}
      </div>

      {postsQuery.data?.items.length ? (
        <div className="max-h-[calc(100vh-17rem)] space-y-5 overflow-y-auto pr-1">
          {postsQuery.data.items.map((post) => (
            <div key={post.id} className={post.imageUrl ? "min-h-[72vh]" : ""}>
              <PostCard post={post} immersive={Boolean(post.imageUrl)} />
            </div>
          ))}
        </div>
      ) : (
        <EmptyState
          title="Пока нет постов"
          description="Создайте первый пост для соседей: объявление, новость, услугу или приглашение на встречу."
          action={
            <Link to="/posts/new" className="inline-flex rounded-2xl bg-ember px-4 py-2.5 text-sm font-semibold text-white">
              Создать пост
            </Link>
          }
        />
      )}
    </div>
  );
}
