export const nullPercent = (
    val: number | null | undefined,
    total: number | null | undefined,
    fixed: number = 2) => {
    return String((((val ?? 0) / (total ?? 1)) * 100).toFixed(fixed)) + "%";
};