import { createContext } from "react";
import { Client } from "./v1/client";
import { IClient } from "./client";

export const ClientContext = createContext<{ queryClient: IClient }>({
	queryClient: new Client(),
});
