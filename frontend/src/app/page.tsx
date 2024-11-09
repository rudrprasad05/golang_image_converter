import Image from "next/image";
import ImageDropBox from "./components/ImageDropBox";
import { Images } from "lucide-react";

export default function Home() {
  return (
    <main className="h-screen w-screen bg-indigo-950">
      <header className="w-full text-indigo-200 px-8 py-5 flex items-center justify-between">
        <div className="flex gap-2 items-center">
          <Images className="text-indigo-200" />
          <h1 className="text-3xl">Image Convertor</h1>
        </div>

        <nav>
          <div className="flex items-center gap-3">
            <h3>Pricing</h3>
            <h3>PDF</h3>
            <h3>Image</h3>
            <h3>Login</h3>
          </div>
        </nav>
      </header>

      <ImageDropBox />
    </main>
  );
}
