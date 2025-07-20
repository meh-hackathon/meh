import type { RouteSectionProps } from "@solidjs/router";
import Navbar from "./components/Navigation";
import NotificationContainer from "./components/Notification";

function App(props:  RouteSectionProps<unknown>) {
	return (
		<div class="flex flex-col items-center min-h-screen bg-gray-100 w-screen">
            <Navbar/>

            {props.children}

            <NotificationContainer />
		</div>
	);
}

export default App;
