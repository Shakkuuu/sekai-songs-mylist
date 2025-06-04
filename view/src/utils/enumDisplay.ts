import { DifficultyType, MusicVideoType } from "../gen/enums/master_pb";
import { ClearType } from "../gen/enums/mylist_pb";

// DifficultyType の表示名マッピング
export const getDifficultyTypeDisplayName = (type: DifficultyType): string => {
  switch (type) {
    case DifficultyType.EASY:
      return "Easy";
    case DifficultyType.NORMAL:
      return "Normal";
    case DifficultyType.HARD:
      return "Hard";
    case DifficultyType.EXPERT:
      return "Expert";
    case DifficultyType.MASTER:
      return "Master";
    case DifficultyType.APPEND:
      return "Append";
    default:
      return "Unspecified";
  }
};

// MusicVideoType の表示名マッピング
export const getMusicVideoTypeDisplayName = (type: MusicVideoType): string => {
  switch (type) {
    case MusicVideoType.MUSIC_VIDEO_TYPE_3D:
      return "3DMV";
    case MusicVideoType.MUSIC_VIDEO_TYPE_2D:
      return "2DMV";
    case MusicVideoType.MUSIC_VIDEO_TYPE_ORIGINAL:
      return "原曲MV";
    default:
      return "未指定";
  }
};

// ClearTypeの表示名マッピング
export const getClearTypeDisplayName = (type: ClearType | undefined): string => {
  switch (type) {
    case ClearType.UNSPECIFIED:
      return "未設定";
    case ClearType.NOT_CLEARED:
      return "未クリア";
    case ClearType.CLEARED:
      return "クリア";
    case ClearType.FULL_COMBO:
      return "フルコンボ";
    case ClearType.ALL_PERFECT:
      return "オールパーフェクト";
    default:
      return "未設定";
  }
};
