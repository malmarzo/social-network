"use client";

import { React } from "react";
import HelloSender from "./components/Buttons/HelloSender";

export default function Home() {

  
  return (
    <div className="container text-white bg-black min-h-screen flex items-center justify-center">
      <div>Hello world </div>
      <HelloSender />
    </div>
  );
}
