"use client";
import Image from "next/image";
import { invokeAPI } from "@/utils/invokeAPI";
import { useRouter } from "next/navigation";
import { React } from "react";
import { useEffect, useState } from "react";
import { AuthProvider } from "@/context/AuthContext";
import CreateGroup from "./createGroup/page"; //
import AuthButton from "./components/Buttons/AuthButtons";
import HelloSender from "./components/Buttons/HelloSender";
export default function Home() {

  
  return (
      <div className="container text-white bg-black min-h-screen flex items-center justify-center">
        Hello world
        {/* <CreateGroup /> */}
        {/* <AuthButton text="create group" href="/createGroup" /> */}
        <AuthButton text="group page " href="/groupPage" />
        <HelloSender />
      </div>
  );
}
