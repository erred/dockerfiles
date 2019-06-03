FROM gcr.io/com-seankhliao/archlinux

RUN useradd -m -u 1000 user
RUN pacman -Syu --noconfirm --needed base base-devel

USER 1000
WORKDIR /workspace
ENTRYPOINT ["/bin/makepkg"]
