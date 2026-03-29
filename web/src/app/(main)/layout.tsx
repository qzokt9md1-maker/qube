import { Sidebar } from "@/components/Sidebar";
import { RightSidebar } from "@/components/RightSidebar";

export default function MainLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex justify-center">
      <Sidebar />
      <main className="xl:ml-[275px] lg:ml-[88px] xl:mr-[350px] w-full max-w-[600px] min-h-screen border-x border-[var(--qube-border)]">
        {children}
      </main>
      <RightSidebar />
    </div>
  );
}
