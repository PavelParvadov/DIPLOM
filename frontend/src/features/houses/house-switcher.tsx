import type { House } from "@/entities/house/model";
import { Select } from "@/shared/ui/primitives";

export function HouseSwitcher({
  houses,
  selectedHouseId,
  onChange,
}: {
  houses: House[];
  selectedHouseId: number | null;
  onChange: (houseId: number) => void;
}) {
  return (
    <Select value={selectedHouseId ?? ""} onChange={(event) => onChange(Number(event.target.value))}>
      <option value="" disabled>
        Выберите дом
      </option>
      {houses.map((house) => (
        <option key={house.id} value={house.id}>
          {house.name} · {house.role === "admin" ? "admin" : "resident"}
        </option>
      ))}
    </Select>
  );
}
