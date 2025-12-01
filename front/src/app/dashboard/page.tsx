import CalendarView from "./components/CalendarView"
import Sidebar from "./components/Sidebar"

export default function Dashboard() {
  return (
    <div className="grid grid-cols-[2fr_6fr] h-screen">
      <Sidebar></Sidebar>
      <CalendarView></CalendarView>
    </div>
  )
} 
