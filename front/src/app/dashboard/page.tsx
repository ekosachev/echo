import CalendarView from "./components/CalendarView"
import Sidebar from "./components/Sidebar"

export default function Dashboard() {



  return (
    <div className="grid grid-cols-1 grid-rows-[1fr_2fr] md:grid-rows-1 md:grid-cols-[1fr_3fr] lg:grid-cols-[1fr_5fr] gap-2 h-screen">

      <Sidebar ></Sidebar>
      <CalendarView></CalendarView>
    </div>
  )
} 
