import { Breadcrumb } from "antd";
import { Link, withRouter } from "react-router-dom";

const breadcrumbMap: Record<string, string> = {
	"/application": "Services",
	"/service-map": "Service Map",
	"/traces": "Traces",
	"/explore": "Explore",
};

export const ShowBreadcrumbs = withRouter((props) => {
	const { location } = props;

	const pathSnippets = location.pathname.split("/").filter((ele) => ele);

	const extraBreadcrumbs = pathSnippets.map((_, idx) => {
		const url = `/${pathSnippets.slice(0, idx + 1).join("/")}`;

		if (!breadcrumbMap[url]) {
			return (
				<Breadcrumb.Item key={url}>
					<Link to={url}>{url.split("/").slice(-1)[0]}</Link>
				</Breadcrumb.Item>
			);
		} else {
			return (
				<Breadcrumb.Item key={url}>
					<Link to={url}>{breadcrumbMap[url]}</Link>
				</Breadcrumb.Item>
			);
		}
	});

	const breadcrumbs = [
		<Breadcrumb.Item key="home">
			<Link to="/"></Link>
		</Breadcrumb.Item>,
	].concat(extraBreadcrumbs);

	return (
		<div>
			<Breadcrumb>{breadcrumbs}</Breadcrumb>
		</div>
	);
});
