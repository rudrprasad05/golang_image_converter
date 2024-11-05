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
import { Loader2, Upload, X } from "lucide-react";
import Image from "next/image";
import { Button } from "@/components/ui/button";

const ImageDropBox = () => {
  const [file, setFile] = useState<File>();
  const router = useRouter();
  const { register, handleSubmit } = useForm();

  const [loadingImage, setloadingImage] = useState<boolean>(false);
  const [formReadyToUpload, setFormReadyToUpload] = useState<boolean>(false);
  const [imageUpload, setImageUpload] = useState(false);
  const [imageUrl, setImageUrl] = useState<string[]>();
  const [isImageInCloud, setIsImageInCloud] = useState(false);

  const form = useForm<UploadImageType>({
    resolver: zodResolver(UploadImageForm),
  });

  const handleImageUpload = async (file: File) => {
    const salt = Date.now();
    console.log(salt.toString() + file.name);
    setImageUpload(true);

    try {
      let data = new FormData();
      data.append("file", file, "image" + salt.toString());

      const res = await fetch("/api/s3-upload", {
        method: "POST",
        body: data,
      })
        .then(() => {
          setImageUpload(false);
          setIsImageInCloud(true);
          //   toast.success("Image Uploaded to Cloud");
        })
        .catch((e) => {
          //   toast("Something went wrong", { description: "Contact site admin" });
        });
      // handle the error
    } catch (e: any) {
      // Handle errors here
      console.error(e);
    }
  };

  function onSubmit(data: UploadImageType) {
    data.file = file;

    const formData = new FormData();
    formData.append("file", data.file);

    axios
      .post(`http://localhost:8080/upload`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
      .then((res) => {
        if (res.status == 200) {
          //   toast.success("Product Created Successfully");
          router.refresh();
        }
      })
      .catch((error) => {
        // toast("Something went wrong", { description: "Contact site admin" });
        console.log("PRODUCT NEW - NewTagButton.tsx", error);
      });
  }
  return (
    <div>
      <main>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-8 w-11/12">
          <input
            type="file"
            onChange={(e) => setFile(e.target.files?.[0] as File)}
          />
          <Button type="submit">Submit</Button>
        </form>
      </main>
    </div>
  );
};

export default ImageDropBox;
