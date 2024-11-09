import React, { useCallback, useState } from "react";
import { DropzoneOptions, useDropzone } from "react-dropzone";

function ImageDropzone({
  image,
  setImage,
}: {
  image?: File;
  setImage: (f: any) => void;
}) {
  // Callback to handle the image drop
  const onDrop = useCallback<NonNullable<DropzoneOptions["onDrop"]>>(
    (acceptedFiles) => {
      const file = acceptedFiles[0];
      setImage(file);
    },
    [setImage]
  );

  // Set up dropzone with the 'onDrop' callback and accept only image files
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    multiple: false,
  });

  return (
    <div
      {...getRootProps()}
      className="text-center w-4/5 h-[150px] border border-dashed mx-auto rounded-sm flex items-center justify-center text-indigo-500 cursor-pointer border-indigo-500"
    >
      <input {...getInputProps()} />
      {isDragActive ? (
        <p>Drop the image here...</p>
      ) : (
        <p>Drag & drop an image here, or click to select one</p>
      )}
      {/* {image && <img src={image} alt="Preview" style={styles.imagePreview} />} */}
    </div>
  );
}

// Define some basic styles

export default ImageDropzone;
