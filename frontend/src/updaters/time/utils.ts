import moment from "moment";

export const Now = () => {
	return Date.now();
};

export const StartTime = (value: number) => {
	return moment(value).format("YYYY-MM-DD hh:mm:ss");
};
