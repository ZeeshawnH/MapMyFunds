export interface PopupState {
  x: number;
  y: number;
  stateCode: string;
  stateName: string;
  visible: boolean;
}

export type TooltipData = {
  x: number;
  y: number;
  name: string;
  stateCode: string;
} | null;
