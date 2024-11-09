"use client";

import React, { useState } from "react";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { FieldValue, useForm } from "react-hook-form";
import { UploadImageForm, UploadImageType } from "@/app/types";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import axios from "axios";
import { toast } from "sonner";
import { ImageIcon, Loader2, Trash, Upload, X } from "lucide-react";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import Link from "next/link";
import { saveAs } from "file-saver";

const ImageDropBox = () => {
  const [file, setFile] = useState<File>();
  const [type, setType] = useState<string>("png");

  const router = useRouter();
  const { register, handleSubmit } = useForm();

  const [loadingImage, setloadingImage] = useState<boolean>(false);
  const [imageUpload, setImageUpload] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>();
  const [isImageInCloud, setIsImageInCloud] = useState(false);

  const handleImageUpload = async (file: File) => {
    setImageUpload(true);
    setloadingImage(true);

    try {
      const formData = new FormData();
      formData.append("file", file);

      axios
        .post(`http://localhost:8080/upload`, formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        })
        .then((res) => {
          setImageUrl(res.data.url);
          setImageUpload(true);
          setloadingImage(false);
          router.refresh();
        })
        .catch((error) => {
          // toast("Something went wrong", { description: "Contact site admin" });
          console.log("ImageDropBox", error);
        });
    } catch (e: any) {
      // Handle errors here
      console.error(e);
    }
  };

  function onSubmit() {
    setloadingImage(true);
    const formData = new FormData();

    if (file && type) {
      formData.append("file", file);
      formData.append("type", type);
    } else {
      console.error("File is undefined");
      console.log(file);
      console.log(type);
      setloadingImage(false);
      return;
    }

    axios
      .post(`http://localhost:8080/convert`, formData)
      .then((res) => {
        if (res.status == 200) {
          console.log(res.data);
          setImageUrl(res.data);
          setloadingImage(false);
        }
      })
      .catch((error) => {
        console.log("ImageDropBox", error);
      });
  }
  const download = () => {
    if (!imageUrl) {
      return;
    }

    fetch(`http://localhost:8080/download?file=${imageUrl}`) // Replace with your actual image URL
      .then((response) => response.blob())
      .then((blob) => {
        saveAs(blob, "image.jpg"); // Set desired filename
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
          <input
            type="file"
            // onChange={(e) => handleImageUpload(e.target.files?.[0] as File)}
            onChange={(e) => setFile(e.target.files?.[0] as File)}
          />
          <section className="w-4/5 mx-auto ">
            {file && (
              <div className="border px-5 py-3 rounded-sm flex gap-8 items-center">
                <div>
                  <ImageIcon />
                </div>
                <div>{file.name}</div>
                <div className="ml-auto">
                  <Trash />
                </div>
              </div>
            )}
          </section>

          <Select onValueChange={(e) => setType(e)} defaultValue={type}>
            <SelectTrigger>
              <SelectValue placeholder="Select a verified email to display" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="png">PNG</SelectItem>
              <SelectItem value="jpeg">JPEG</SelectItem>
              <SelectItem value="jpg">JPG</SelectItem>
              <SelectItem value="webp">WEBP</SelectItem>
            </SelectContent>
          </Select>

          <div className="flex items-center gap-3">
            <Button type="submit">Submit</Button>
            {loadingImage && (
              <div>
                <Loader2 className="animate-spin" />
              </div>
            )}
          </div>
        </form>
      </main>
      {imageUrl && (
        <section>
          <div className="w-[100px] h-[100px]">
            <Image
              src={imageUrl}
              alt="image"
              height={50}
              width={50}
              className="w-auto h-auto object-cover aspect-square rounded-sm"
            />
          </div>
          <div>
            <Button onClick={download}>Download!</Button>
          </div>
        </section>
      )}
    </div>
  );
};

export default ImageDropBox;
