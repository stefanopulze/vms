

export function customJsonReceiver(json: string): any {
  return JSON.parse(json, dateReviver);
}

function dateReviver(key: string, value: any): any {
  if ("timestamp" === key) {
    return isEmptyDate(value) ? null : new Date(value);
  }

  return value;
}

function isEmptyDate(date: string | null) {
  return date === null || date === "0001-01-01T00:00:00Z"
}
