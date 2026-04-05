import type { Comment } from "@/entities/comment/model";
import { formatDate } from "@/shared/lib/format";
import { classNames } from "@/shared/lib/format";

export function CommentList({ comments }: { comments: Comment[] }) {
  return (
    <div className="space-y-3">
      {comments.map((comment, index) => (
        <article
          key={comment.id}
          className={classNames(
            "rounded-[24px] border px-4 py-4",
            index % 2 === 0 ? "border-slate-200 bg-white" : "border-sand/70 bg-sand/20",
          )}
        >
          <div className="mb-2 flex items-center justify-between gap-4 text-sm">
            <span className="font-semibold text-ink">{comment.authorName}</span>
            <span className="text-slate-500">{formatDate(comment.createdAt)}</span>
          </div>
          <p className="whitespace-pre-wrap text-sm leading-6 text-slate-700">{comment.content}</p>
        </article>
      ))}
    </div>
  );
}
