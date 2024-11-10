"use client";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import axios from "axios";
import clsx from "clsx";
import { saveAs } from "file-saver";
import { Download, ImageIcon, RefreshCcw, X } from "lucide-react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";
import ImageDropzone from "./ImageDropzone";

const ImageDropBox = () => {
  const [file, setFile] = useState<File>();
  const [type, setType] = useState<string>("png");

  const [loadingImage, setloadingImage] = useState<boolean>(false);
  const [imageUrl, setImageUrl] = useState<string>();

  function onSubmit() {
    setloadingImage(true);
    const formData = new FormData();

    if (file && type) {
      formData.append("file", file);
      formData.append("type", type);
    } else {
      toast.error("Image Uploaded to Cloud");
      console.error("File is undefined");
      console.log(file);
      console.log(type);
      setloadingImage(false);
      return;
    }

    axios
      .post(`https://44.204.76.51:8080/convert`, formData)
      .then((res) => {
        if (res.status == 200) {
          console.log(res.data);
          setImageUrl(res.data);
          setloadingImage(false);
        }
      })
      .catch((error) => {
        console.log("ImageDropBox", error);
        setloadingImage(false);
      });
  }
  const download = () => {
    if (!imageUrl) {
      return;
    }

    fetch(`https://44.204.76.51:8080/download?file=${imageUrl}`) // Replace with your actual image URL
      .then((response) => response.blob())
      .then((blob) => {
        saveAs(blob, "converted_image." + type); // Set desired filename
      })
      .catch((error) => console.error("Error downloading the image:", error));
  };
  return (
    <div>
      <main>
        <form
          onSubmit={(e) => {
            e.preventDefault(); // Prevent the default form submission
            e.stopPropagation();
            onSubmit();
          }}
          className="space-y-5 pt-8"
        >
          <ImageDropzone setImage={setFile} />

          <section className="w-4/5 mx-auto flex items-center gap-4">
            <section className="grow">
              <Select onValueChange={(e) => setType(e)} defaultValue={type}>
                <SelectTrigger className="bg-indigo-700 border-none">
                  <SelectValue placeholder="Select a verified email to display" />
                </SelectTrigger>
                <SelectContent className="bg-indigo-700 border-none">
                  <SelectItem className="border-none" value="png">
                    PNG
                  </SelectItem>
                  <SelectItem className="border-none" value="jpeg">
                    JPEG
                  </SelectItem>
                  <SelectItem className="border-none" value="jpg">
                    JPG
                  </SelectItem>
                  <SelectItem className="border-none" value="webp">
                    WEBP
                  </SelectItem>
                </SelectContent>
              </Select>
            </section>

            <div className="flex items-center gap-3">
              <Button
                type="submit"
                className="flex items-center gap-2 bg-indigo-200"
              >
                <RefreshCcw
                  className={clsx(loadingImage ? "animate-spin" : "")}
                />
                Convert
              </Button>
            </div>
          </section>

          <section className="w-4/5 mx-auto ">
            {file && (
              <div className="bg-indigo-700 px-5 py-3 rounded-sm flex gap-8 items-center">
                <div>
                  <ImageIcon className="text-indigo-300" />
                </div>
                <div className="truncate">{file.name}</div>
                <div className="text-indigo-300">
                  {file.type.split("/").pop()}
                </div>
                <div className="ml-auto">
                  <div
                    onClick={() => setFile(undefined)}
                    className="bg-destructive cursor-pointer p-[2px] rounded-full w-5 h-5 flex items-center justify-center"
                  >
                    <X className="text-indigo-200" />
                  </div>
                </div>
              </div>
            )}
          </section>
        </form>
      </main>
      <section className="w-4/5 mx-auto mt-6">
        {imageUrl && (
          <section className="bg-indigo-700 rounded-sm flex justify-between items-center">
            <div className="w-[100px] h-[100px]">
              <Image
                src={imageUrl}
                alt="image"
                height={100}
                width={100}
                className="w-auto h-auto object-cover aspect-square rounded-l-sm"
              />
            </div>

            <div className="pr-6">
              <Button size={"sm"} onClick={download}>
                <Download size={10} />
                Download
              </Button>
            </div>
          </section>
        )}
      </section>
    </div>
  );
};

export default ImageDropBox;
