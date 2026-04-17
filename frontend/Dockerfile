FROM node:20-slim
WORKDIR /app
COPY package.json ./
# Glibc base (not alpine/musl) so Next's prebuilt @next/swc-linux-*-gnu
# binary is ABI-compatible.
#
# We deliberately do NOT copy package-lock.json into the image. The
# committed lockfile is generated on macOS and pins the darwin-arm64 SWC
# binary as the "resolved" optional dep; both `npm ci` and `npm install`
# will honour the lockfile's resolution and skip installing
# @next/swc-linux-arm64-gnu for this linux build, leaving `next build` to
# try to fetch SWC at runtime — which fails inside sandboxed builders.
# Resolving from package.json alone forces npm to pick the right native
# binary for the current (linux) build platform.
RUN npm install --include=optional --no-audit --no-fund
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
