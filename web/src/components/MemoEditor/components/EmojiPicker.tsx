import EmojiPicker, { type EmojiClickData, Theme } from "emoji-picker-react";
import { SmileIcon } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { useTranslate } from "@/utils/i18n";

interface EmojiPickerButtonProps {
  onEmojiSelect: (emoji: string) => void;
  disabled?: boolean;
}

export const EmojiPickerButton = ({ onEmojiSelect, disabled }: EmojiPickerButtonProps) => {
  const t = useTranslate();
  const [open, setOpen] = useState(false);

  const handleEmojiClick = (emojiData: EmojiClickData) => {
    onEmojiSelect(emojiData.emoji);
    setOpen(false);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="ghost" size="sm" disabled={disabled} title={t("tooltip.insert-emoji")}>
          <SmileIcon className="w-4 h-4" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[350px] p-0 border-0" align="start">
        <EmojiPicker
          onEmojiClick={handleEmojiClick}
          theme={Theme.AUTO}
          width="100%"
          height={400}
          searchPlaceHolder={t("editor.search-emoji")}
          previewConfig={{ showPreview: false }}
        />
      </PopoverContent>
    </Popover>
  );
};
