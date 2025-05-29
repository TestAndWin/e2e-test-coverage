#!/bin/bash

echo "📦 Install JS libs (ui/)..."
cd ui
npm ci  
cd ..

echo "🐹 Install golang libs (api/)..."
cd api
go mod tidy 
cd ..

echo "✅ Setup done."
