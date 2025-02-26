
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/layout/Layout";
import Dashboard from "./pages/Dashboard";
import StudyActivities from "./pages/StudyActivities";
import StudyActivityDetails from "./pages/StudyActivityDetails";
import StudyActivityLauncher from "./pages/StudyActivityLauncher";
import Words from "./pages/Words";
import WordGroups from "./pages/WordGroups";
import WordGroup from "./pages/WordGroup";
import Sessions from "./pages/Sessions";
import Settings from "./pages/Settings";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route element={<Layout />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/study-activities" element={<StudyActivities />} />
            <Route path="/study-activities/:id" element={<StudyActivityDetails />} />
            <Route path="/study-activities/:id/launch" element={<StudyActivityLauncher />} />
            <Route path="/words" element={<Words />} />
            <Route path="/groups" element={<WordGroups />} />
            <Route path="/groups/:id" element={<WordGroup />} />
            <Route path="/sessions" element={<Sessions />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
