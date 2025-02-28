"use client";
import Image from "next/image";
import { invokeAPI } from "@/utils/invokeAPI";
import { useRouter } from "next/navigation";
import { React } from "react";
import { useEffect, useState } from "react";
import { AuthProvider } from "@/context/AuthContext";

export default function Home() {

  
  return (
      <div className="container text-white bg-black min-h-screen flex items-center justify-center">
        Hello world
      </div>
  );
}
