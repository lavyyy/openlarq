import type { MouseEventHandler, FocusEventHandler } from 'svelte/elements';

export const getLastMonday = (start: Date) => {
	const diff = (start.getDay() + 6) % 7;
	start.setDate(start.getDate() - diff);
	return start;
};

export const getColor = (
	colors: {
		light: string[];
		dark: string[];
	},
	max: number,
	value: number,
	siteTheme: 'dark' | 'light'
) => {
	const themeColors = colors[siteTheme === 'dark' ? 'dark' : 'light'];

	console.log(value);
	if (!value) return themeColors[0];
	const p = (value / max) * (themeColors.length - 1);
	return themeColors[Math.floor(p)];
};

export const getCalendar = (data: { [key: string]: number }, year: number) => {
	const base = getLastMonday(new Date(year, 0, 1));
	const out: {
		max: number;
		calendar: ({ date: string; value: number } | undefined)[][];
		monthLabels?: string[];
	} = {
		max: 0,
		calendar: []
	};

	out.max = 0;
	out.calendar = Array.from({ length: 7 }, (_, i) => {
		const start = new Date(base);
		start.setDate(start.getDate() + i);
		return Array.from({ length: 53 }, (_, j) => {
			const day = new Date(start);
			day.setDate(start.getDate() + j * 7);
			if (day.getFullYear() == year) {
				const date = day.toISOString().split('T')[0];
				const value = data[date] ?? 0;
				if (value > out?.max) {
					out.max = value;
				}
				return { date, value };
			}
		});
	});

	const monthLabels: string[] = [];
	let lastMonth = -1;
	for (let week = 0; week < 53; week++) {
		let label = '';
		for (let day = 0; day < 7; day++) {
			const cell = out.calendar[day][week];
			if (cell && cell.date) {
				const dateObj = new Date(cell.date);
				const month = dateObj.getMonth();
				const cellYear = dateObj.getFullYear();
				if (cellYear === year && month !== lastMonth) {
					label = dateObj.toLocaleString('default', { month: 'short' });
					lastMonth = month;
				}
				break;
			}
		}
		monthLabels.push(label);
	}
	out.monthLabels = monthLabels;

	return out;
};

export const calMonths = [
	'Jan',
	'Feb',
	'Mar',
	'Apr',
	'May',
	'Jun',
	'Jul',
	'Aug',
	'Sep',
	'Oct',
	'Nov',
	'Dec'
];

export const getColSpanForMonth = (month: string) => {
	switch (month) {
		case 'Jan':
			return 5;
		case 'Feb':
			return 4;
		case 'Mar':
			return 4;
		case 'Apr':
			return 5;
		case 'May':
			return 4;
		case 'Jun':
			return 4;
		case 'Jul':
			return 5;
		case 'Aug':
			return 4;
		case 'Sep':
			return 4;
		case 'Oct':
			return 5;
		case 'Nov':
			return 4;
		case 'Dec':
			return 4;
	}
};

export type Props = {
	data: { [key: string]: number };

	year?: number;
	lday?: boolean;
	lmonth?: boolean;
	colors?: {
		light: string[];
		dark: string[];
	};
	className?: string;

	onclick?: MouseEventHandler<HTMLTableCellElement>;
	onmouseout?: MouseEventHandler<HTMLTableCellElement>;
	onmouseover?: MouseEventHandler<HTMLTableCellElement>;
	onfocus?: FocusEventHandler<HTMLTableCellElement>;
	onblur?: FocusEventHandler<HTMLTableCellElement>;
};
