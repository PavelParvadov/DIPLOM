import { NavLink } from "react-router-dom";

import { classNames } from "@/shared/lib/format";

const items = [
  { to: "/admin/categories", label: "Категории" },
  { to: "/admin/posts", label: "Посты" },
  { to: "/admin/invites", label: "Коды" },
];

export function AdminNav() {
  return (
    <div className="flex flex-wrap gap-3">
      {items.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          className={({ isActive }) =>
            classNames(
              "rounded-2xl border px-4 py-2 text-sm font-semibold transition",
              isActive ? "border-moss bg-moss text-white" : "border-slate-200 bg-white text-ink hover:border-moss hover:text-moss",
            )
          }
        >
          {item.label}
        </NavLink>
      ))}
    </div>
  );
}
