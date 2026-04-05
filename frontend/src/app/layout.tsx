import { useEffect } from "react";
import { Link, NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";

import { useSession } from "@/app/providers";
import { HouseSwitcher } from "@/features/houses/house-switcher";
import { useHouses } from "@/features/houses/queries";
import { classNames } from "@/shared/lib/format";
import { Button, SecondaryButton } from "@/shared/ui/primitives";

export function AppLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, selectedHouseId, setSelectedHouseId } = useSession();
  const housesQuery = useHouses(Boolean(user));

  useEffect(() => {
    if (!housesQuery.data?.length) {
      return;
    }

    const exists = housesQuery.data.some((house) => house.id === selectedHouseId);
    if (!selectedHouseId || !exists) {
      setSelectedHouseId(housesQuery.data[0].id);
    }
  }, [housesQuery.data, selectedHouseId, setSelectedHouseId]);

  const selectedHouse = housesQuery.data?.find((house) => house.id === selectedHouseId) ?? null;
  const isAdmin = selectedHouse?.role === "admin";
  const navItems = [
    { to: "/posts", label: "Лента" },
    { to: "/chat", label: "Чат" },
    { to: "/posts/new", label: "Новый пост" },
    { to: "/join", label: "Дома и коды" },
    ...(isAdmin ? [{ to: "/admin/categories", label: "Админ" }] : []),
  ];

  return (
    <div className="min-h-screen px-4 py-6 md:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl space-y-6">
        <header className="rounded-[32px] border border-white/70 bg-white/85 px-6 py-5 shadow-soft backdrop-blur">
          <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
            <div className="space-y-1">
              <Link to="/posts" className="font-display text-3xl text-ink">
                HappyHouse
              </Link>
              <p className="text-sm text-slate-500">Цифровой клуб соседей для общения, новостей дома и общих обсуждений в одном аккуратном интерфейсе.</p>
            </div>
            <div className="flex flex-col gap-4 lg:items-end">
              <div className="flex flex-wrap items-center gap-3">
                <span className="rounded-full bg-sand px-3 py-2 text-sm font-medium text-ink">{user?.displayName}</span>
                <SecondaryButton
                  onClick={() => {
                    logout();
                    navigate("/login");
                  }}
                >
                  Выйти
                </SecondaryButton>
              </div>
              <div className="min-w-[260px]">
                {housesQuery.data?.length ? (
                  <HouseSwitcher houses={housesQuery.data} selectedHouseId={selectedHouseId} onChange={setSelectedHouseId} />
                ) : (
                  <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 px-4 py-3 text-sm text-slate-500">
                    У вас пока нет домов. Создайте новый дом или присоединитесь по invite code.
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="mt-5 flex flex-wrap gap-3">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  classNames(
                    "rounded-full px-4 py-2 text-sm font-semibold transition",
                    isActive ? "bg-ink text-white" : "bg-slate-100 text-slate-600 hover:bg-slate-200",
                  )
                }
              >
                {item.label}
              </NavLink>
            ))}
          </div>
        </header>

        {selectedHouse ? (
          <section className="rounded-[28px] bg-sand/60 px-5 py-4 text-sm text-slate-700 shadow-soft">
            Активный дом: <span className="font-semibold text-ink">{selectedHouse.name}</span> · {selectedHouse.address}
          </section>
        ) : null}

        {housesQuery.isError ? (
          <section className="rounded-[28px] bg-red-50 px-5 py-4 text-sm text-red-700">Не удалось загрузить список домов.</section>
        ) : null}

        {!selectedHouse && location.pathname !== "/join" ? (
          <section className="rounded-[28px] border border-dashed border-slate-300 bg-white px-6 py-8 text-center shadow-soft">
            <h2 className="text-xl font-semibold text-ink">Сначала создайте дом или присоединитесь к существующему</h2>
            <p className="mt-2 text-sm text-slate-600">Без активного дома нельзя открыть ленту, чат, посты и инструменты администрирования.</p>
            <Button className="mt-4" onClick={() => navigate("/join")}>
              Перейти к управлению домами
            </Button>
          </section>
        ) : (
          <Outlet context={{ selectedHouse, houses: housesQuery.data ?? [] }} />
        )}
      </div>
    </div>
  );
}
