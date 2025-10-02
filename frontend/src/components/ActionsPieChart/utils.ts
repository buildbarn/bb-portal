const pieChartColors = [
  "#E6194B",
  "#3CB44B",
  "#D8BD14",
  "#4363D8",
  "#F58231",
  "#911EB4",
  "#58D1C9",
  "#F032E6",
  "#A9D134",
  "#F37FB7",
  "#008080",
  "#E0529C",
  "#9A6324",
  "#F3B765",
  "#D32029",
  "#8FD460",
  "#808000",
  "#AB7AE0",
  "#65A9F3",
  "#808080",
];

export const chartColor = (index: number): string => {
  return pieChartColors[index % pieChartColors.length];
};
