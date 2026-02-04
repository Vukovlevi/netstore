import type { DeletedAt, NullTime } from "./DeletedAt";

type WeekDay = {
  id: number;
  name: string;
};

type OpenHour = {
  id: number;
  opensAt: string;
  closesAt: string;
  specialDate: NullTime;
  weekDayIds: number[];
  weekDays: string[];
  deletedAt: DeletedAt | null;
};

class OpenHourClass {
  id = 0;
  opensAt = "";
  closesAt = "";
  specialDate: NullTime = { Valid: false, Time: "0001-01-01T00:00:00Z" };
  weekDayIds: number[] = [];
  weekDays: string[] = [];
  deletedAt: DeletedAt | null = null;
  inputSpecialDate = "";
  wasWeekDayChanged = false;

  constructor(openHour: OpenHour | null) {
    if (openHour == null) {
      return;
    }

    this.id = openHour.id;
    this.opensAt = openHour.opensAt.substring(0, openHour.opensAt.length - 3);
    this.closesAt = openHour.closesAt.substring(
      0,
      openHour.closesAt.length - 3
    );
    if (openHour.specialDate && openHour.specialDate.Valid) {
      this.specialDate = openHour.specialDate;
      this.inputSpecialDate = new Date(openHour.specialDate.Time.split("T")[0]!)
        .toISOString()
        .substring(0, 10);
    }
    if (openHour.weekDayIds) this.weekDayIds = openHour.weekDayIds;
    if (openHour.weekDays) this.weekDays = openHour.weekDays;
    this.deletedAt = openHour.deletedAt;
  }

  toOpenHour(): OpenHour {
    return {
      id: this.id,
      opensAt: this.opensAt + ":00",
      closesAt: this.closesAt + ":00",
      specialDate:
        this.inputSpecialDate == ""
          ? { Valid: false, Time: "0001-01-01T00:00:00Z" }
          : {
              Valid: true,
              Time: new Date(this.inputSpecialDate).toISOString(),
            },
      weekDayIds: this.weekDayIds,
      weekDays: this.weekDays,
      deletedAt: this.deletedAt,
    };
  }

  compare(openHour: OpenHour): boolean {
    if (this.inputSpecialDate == "" && openHour.specialDate.Valid) return false;
    if (this.inputSpecialDate != "" && !openHour.specialDate.Valid)
      return false;
    if (
      this.inputSpecialDate != "" &&
      openHour.specialDate.Valid &&
      openHour.specialDate.Time != new Date(this.inputSpecialDate).toISOString()
    )
      return false;
    return (
      openHour.opensAt == this.opensAt + ":00" &&
      openHour.closesAt == this.closesAt + ":00" &&
      !this.wasWeekDayChanged
    );
  }
}

export type { OpenHour, WeekDay };
export { OpenHourClass };
