"use client"

import { useState } from "react"

export default function UploadTestPage() {
  const [orderId, setOrderId] = useState("")
  const [files, setFiles] = useState<FileList | null>(null)
  const [message, setMessage] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!orderId || !files || files.length === 0) {
      setMessage("Please provide an order id & at least 1 files")
      return
    }

    //! use form data instead JSON.stringify()
    const formData = new FormData()
    formData.append("order_id", orderId)

    //! loop through all selected file and append to tech_packs
    for (let i = 0; i < files.length; i++) {
      formData.append("tech_packs", files[i])
    }

    try {
      //! sending POST req
      //! don't set Content-type header
      const res = await fetch("http://localhost:8080/orders/images", {
        method: "POST",
        body: formData,
      })

      const responseText = await res.text()

      if (res.ok) {
        setMessage("success " + responseText)
      } else {
        setMessage("failed error: " + responseText)
      }
    } catch (error) {
      setMessage("Network error occured")
    }
  }

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-6 text-gray-800">Upload Tech Packs</h1>
        
        <form onSubmit={handleSubmit} encType="multipart/form-data" className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Order ID</label>
            <input
              type="number"
              value={orderId}
              onChange={(e) => setOrderId(e.target.value)}
              className="w-full border border-gray-300 rounded p-2 text-gray-800"
              placeholder="e.g., 2"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Select Images</label>
            {/* The 'multiple' attribute allows selecting more than one file */}
            <input
              type="file"
              multiple
              accept="image/*"
              onChange={(e) => setFiles(e.target.files)}
              className="w-full text-gray-700"
            />
          </div>

          <button
            type="submit"
            className="w-full bg-blue-600 text-white font-bold py-2 px-4 rounded hover:bg-blue-700 transition"
          >
            Upload Images
          </button>
        </form>

        {message && (
          <div className="mt-6 p-4 bg-gray-100 rounded text-gray-800 font-medium border border-gray-200">
            {message}
          </div>
        )}
      </div>
    </main>
  )
}